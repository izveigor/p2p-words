import random
from concurrent import futures

import grpc
import spacy

import lemmatizer_pb2_grpc
from lemmatizer_pb2 import LemmatizedWord, LemmatizerResponse

model = spacy.load("ru_core_news_sm")
MAX = 2**16
ALPHABET = "абвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"
P2P_SVC_URL = "localhost:50051"


class LemmatizerService(lemmatizer_pb2_grpc.LemmatizersServicer):
    def Lemmatize(self, request, context):
        text = request.text
        lemmatized_text = []
        for token in model(text):
            lemma = str(token.lemma_)
            will_lemmatize = True
            for char in lemma:
                if char not in ALPHABET:
                    will_lemmatize = False
                    break

            if will_lemmatize:
                lemmatized_text.append(
                    LemmatizedWord(
                        id=random.randint(1, MAX),
                        word=lemma,
                        sentence=str(token.sent),
                    )
                )

        return LemmatizerResponse(
            words=lemmatized_text,
        )


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    lemmatizer_pb2_grpc.add_LemmatizersServicer_to_server(LemmatizerService(), server)
    server.add_insecure_port(P2P_SVC_URL)
    server.start()
    server.wait_for_termination()


if __name__ == "__main__":
    serve()

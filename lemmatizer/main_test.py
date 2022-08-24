import unittest

from main import LemmatizerService
from lemmatizer_pb2 import LemmatizerRequest


class TestLemmatizer(unittest.TestCase):
    def test_lemmatize(self):
        service = LemmatizerService()
        request = LemmatizerRequest(
            text="Жуки-златки среднего размера. Тело узкое и "
            + "удлинённое, клиновидное."  # https://ru.wikipedia.org/wiki/%D0%AF%D1%81%D0%B5%D0%BD%D0%B5%D0%B2%D0%B0%D1%8F_%D0%B8%D0%B7%D1%83%D0%BC%D1%80%D1%83%D0%B4%D0%BD%D0%B0%D1%8F_%D1%83%D0%B7%D0%BA%D0%BE%D1%82%D0%B5%D0%BB%D0%B0%D1%8F_%D0%B7%D0%BB%D0%B0%D1%82%D0%BA%D0%B0
        )
        response = service.Lemmatize(request, None)

        data_answer = [
            ("жук", "Жуки-златки среднего размера."),
            ("златка", "Жуки-златки среднего размера."),
            ("средний", "Жуки-златки среднего размера."),
            ("размер", "Жуки-златки среднего размера."),
            ("тело", "Тело узкое и удлинённое, клиновидное."),
            ("узкий", "Тело узкое и удлинённое, клиновидное."),
            ("и", "Тело узкое и удлинённое, клиновидное."),
            ("удлинённый", "Тело узкое и удлинённое, клиновидное."),
            ("клиновидный", "Тело узкое и удлинённое, клиновидное."),
        ]

        for data, answer in zip(response.words, data_answer):
            self.assertEqual(data.word, answer[0])
            self.assertEqual(data.sentence, answer[1])


if __name__ == "__main__":
    unittest.main()

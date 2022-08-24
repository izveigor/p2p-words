proto:
	python -m grpc_tools.protoc -I . --python_out=. --grpc_python_out=. ./lemmatizer.proto

start:
	python main.py

test:
	python main_test.py
FROM python:3.6

ENV SERVICE_PATH="disaster-service"

WORKDIR /code

COPY ./${SERVICE_PATH}/requirements.txt .

# install dependencies
RUN pip install -r requirements.txt

COPY ./${SERVICE_PATH}/src/ .

CMD [ "python", "./app.py" ]

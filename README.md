# Fronting Postgres


## Setting up environment

```
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

## running the webserver

```
uvicorn app:app --reload --host 0.0.0.0
```

## running the locust job

```
http://localhost:8099/
```

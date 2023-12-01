from fastapi import FastAPI
from fastapi.encoders import jsonable_encoder
from fastapi.responses import JSONResponse
from typing import Optional
from sqlmodel import Field, SQLModel, Session, create_engine, select
from faker import Faker
import os, datetime, redis, pickle

app = FastAPI()
dburl = os.environ.get('DB_URL', "postgresql://postgres:PgDbFTW15@localhost:5432/profiles")
redis_host = os.environ.get('REDIS_HOST', "localhost")
redis_port = os.environ.get('REDIS_PORT', "6379")
redis_pass = os.environ.get('REDIS_PASS', None)

engine = create_engine(dburl)
r = redis.StrictRedis(
        host=redis_host,
        password=redis_pass,
        port=int(redis_port),
        )

##
class Profile(SQLModel, table=True):
    id: Optional[int] = Field(default=None, primary_key=True)
    first_name: str
    last_name: str
    email: str
    job: str
    ssn: str
    phone: str
    dob: datetime.date = Field(nullable=False)

## Start routes
@app.get("/")
async def root():
    return {"message": "pong"}

@app.get("/profile/{profile_id}")
async def read_profile(profile_id):
    if r.exists('ENABLE_CACHING'):
        data=r.get("user:{}".format(profile_id))
        return JSONResponse(jsonable_encoder(pickle.loads(data)))

    else:
        data = {}
        with Session(engine) as session:
            statement=select(Profile).where(Profile.id == profile_id)
            results = session.exec(statement)
            for result in results:
                data = result
        return JSONResponse(jsonable_encoder(data))

@app.get("/load")
async def loaddata():
    fake = Faker()
    SQLModel.metadata.create_all(engine)
    with Session(engine) as session:
        for x in range(100000):
            fn=fake.first_name()
            ln=fake.last_name()
            dom=fake.domain_name()
            profile = Profile(
                        id=x,
                        first_name=fn,
                        last_name=ln,
                        email="{0}_{1}@{2}".format(fn, ln, dom).lower(),
                        job=fake.job(),
                        ssn=fake.ssn(),
                        phone=fake.phone_number(),
                        dob=fake.date_between(),
                        )
            r.set("user:{}".format(x), pickle.dumps(profile))
            session.add(profile)
            if x % 100 == 0:
                session.commit()
        session.commit()

    return {"message": "data_loaded"}

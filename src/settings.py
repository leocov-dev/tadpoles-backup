import environs
import requests

env = environs.Env()
env.read_env()


class Config:
    OAUTH_TOKEN = env("OAUTH_TOKEN")
    API_URL = "https://www.tadpoles.com/remote/v1"


def get_client():
    rc = requests.Session()
    rc.headers = {"Cookie": f"DgU00={conf.OAUTH_TOKEN}"}
    return rc


conf = Config()
client = get_client()

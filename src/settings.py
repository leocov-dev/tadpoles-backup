import environs

env = environs.Env()
env.read_env()


class Config:
    CLIENT_ID = env("CLIENT_ID")
    CLIENT_SECRET = env("CLIENT_SECRET")


conf = Config()

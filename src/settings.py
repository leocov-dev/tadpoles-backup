import environs

env = environs.Env()
env.read_env()


class Config:
    CLIENT_ID = env("CLIENT_ID")


conf = Config()

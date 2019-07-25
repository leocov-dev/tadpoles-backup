import os

from selenium import webdriver


def selenium_flow():
    """"""

    driver = webdriver.Chrome(os.path.join(os.path.dirname(__file__), "..", "webdriver", "chromedriver.exe"))
    driver.get("https://www.tadpoles.com/home_or_work")
    elem = driver.find_element_by_class_name("other-login-button")

    print(elem.get_attribute("data-bind"))


if __name__ == '__main__':
    selenium_flow()

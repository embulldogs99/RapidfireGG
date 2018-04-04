from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.ui import WebDriverWait
from selenium import webdriver
from selenium.webdriver.common.by import By

email="abcdefg@us.org"


EMAILFIELD = (By.ID, "i0116")
PASSWORDFIELD = (By.ID, "i0118")
NEXTBUTTON = (By.ID, "idSIButton9")

browser = webdriver.Firefox()
browser.get('https://login.live.com/login.srf')

# wait for email field and enter email
WebDriverWait(browser, 10).until(EC.element_to_be_clickable(EMAILFIELD)).send_keys(email)

# Click Next
WebDriverWait(browser, 10).until(EC.element_to_be_clickable(NEXTBUTTON)).click()

# wait for password field and enter password
WebDriverWait(browser, 10).until(EC.element_to_be_clickable(PASSWORDFIELD)).send_keys("Rumble34")

# Click Login - same id?
WebDriverWait(browser, 10).until(EC.element_to_be_clickable(NEXTBUTTON)).click()

print("made it")

export function isPasswordStrong(password) {
  const strongPasswordRegex =
    /^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%^&*])(?=.{8,})/;
  return strongPasswordRegex.test(password);
}
export function validateEmail(email) {
  var re =
    /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
  return re.test(String(email).toLowerCase());
}
export const Url = "http://localhost:8000/";

export const getJson = async (endpoint, aptions) => {
  try {
    const res = await fetch(`${Url}${endpoint}`, aptions);
    if (!res.ok) {
      const error = await res.json();
      throw Error(`${error.error} ${res.statusText} ${res.status} `);
    }
    return res.json();
  } catch (error) {
    throw error;
  }
};


const avarege 
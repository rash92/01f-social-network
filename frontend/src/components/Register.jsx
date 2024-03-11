import FormGroup from "./FormGroup";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import {useState} from "react";
import {validateEmail, getJson} from "../helpers/helpers";

export default function Register({setShowRegisterForm}) {
  const [error, setError] = useState({isError: false, message: ""});

  const [emailValue, setEmailValue] = useState("");
  const [passwordValue, setPasswordValue] = useState("");
  const [firstNameValue, setFirstNameValue] = useState("");
  const [lastNameValue, setLastNameValue] = useState("");
  const [usernameValue, setUsernameValue] = useState("");
  const [dateOfBirthValue, setDateOfBirthValue] = useState("");
  const [avatarValue, setAvatarValue] = useState({value: "", file: null});
  const [aboutMeValue, setAboutMeValue] = useState("");

  const isNotEmty = (val) => {
    return val.length > 0;
  };

  const isPasswordStrong = (val) => {
    return val.length > 7;
  };

  const registerHandler = async (e) => {
    e.preventDefault();

    const formIsValid =
      validateEmail(emailValue) &&
      isPasswordStrong(passwordValue) &&
      isNotEmty(firstNameValue) &&
      isNotEmty(lastNameValue) &&
      isNotEmty(dateOfBirthValue); 

    

    if (  !formIsValid) {
      alert("Please fill all fields");
      return;
    }
    const formData = new FormData();
    formData.append("image", avatarValue.file);
    formData.append("email", emailValue);
    formData.append("password", passwordValue);
    formData.append("firstName", firstNameValue);
    formData.append("lastName", lastNameValue);
    formData.append("nickname", usernameValue);
    formData.append("DOB", dateOfBirthValue);
    formData.append("aboutMe", aboutMeValue);

    try {
      const res = await getJson("newUser", {
        method: "POST",
        body: formData,
      });

      setShowRegisterForm(false);
    } catch (err) {
      setError({isError: true, message: err.message});
    }
  };

  return (
    <Form onSubmit={registerHandler}>
      {error.isError && <div className="invalid-feedback">{error.message}</div>}
      <FormGroup
        value={emailValue}
        setValue={setEmailValue}
        type="email"
        Label="Email"
        required={true}
        Text="(Required.) We'll never share your email with anyone else."
        placeholder="name@example.com"
      />
      <FormGroup
        type="password"
        Label="Password"
        required={true}
        Text="(Required.) Password must be at least 8 characters long."
        value={passwordValue}
        setValue={setPasswordValue}
        placeholder="12345678"
      />
      <FormGroup
        value={firstNameValue}
        setValue={setFirstNameValue}
        required={true}
        type="text"
        Label="First Name"
        Text="(Required.) We value your first name."
        placeholder="John"
      />
      <FormGroup
        value={lastNameValue}
        setValue={setLastNameValue}
        type="text"
        Label="Last Name"
        required={true}
        Text="(Required.) We value your last name more."
        placeholder="Doe"
      />
      <FormGroup
        value={dateOfBirthValue}
        setValue={setDateOfBirthValue}
        type="date"
        required={true}
        Label="Date of Birth"
        Text="(Required.)"
        placeholder="MM/DD/YYYY"
      />
      <FormGroup
        value={usernameValue}
        setValue={setUsernameValue}
        type="text"
        Label="Username"
        Text="We're neutral about your username."
        placeholder="john_doe"
      />
      <FormGroup
        value={avatarValue}
        setValue={setAvatarValue}
        type="file"
        accept="image/*"
        Label="Avatar"
        Text="JPG, PNG, GIF."
      />
      <FormGroup
        value={aboutMeValue}
        setValue={setAboutMeValue}
        type="textarea"
        Label="About Me"
        Text="Tell us about yourself."
        placeholder="I'm a cool guy."
      />
      <Button variant="primary" type="submit">
        Submit
      </Button>
      <div>
        if you already have an account, please login here:
        <Button
          variant="primary"
          onClick={() => {
            setShowRegisterForm(false);
          }}
        >
          login
        </Button>
      </div>
    </Form>
  );
}

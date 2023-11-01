import FormGroup from "./FormGroup";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import { useState } from "react";
import { Container } from "react-bootstrap";

export default function Register({ setShowRegisterForm }) {
  const [emailValue, setEmailValue] = useState("");
  const [passwordValue, setPasswordValue] = useState("");
  const [firstNameValue, setFirstNameValue] = useState("");
  const [lastNameValue, setLastNameValue] = useState("");
  const [usernameValue, setUsernameValue] = useState("");
  const [dateOfBirthValue, setDateOfBirthValue] = useState("");
  const [avatarValue, setAvatarValue] = useState("");
  const [aboutMeValue, setAboutMeValue] = useState("");
  return (
    <Form>
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

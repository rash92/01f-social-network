import FormGroup from "./FormGroup";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import { useState } from "react";

export default function Login({ setShowRegisterForm }) {
  const [emailValue, setEmailValue] = useState("");
  const [passwordValue, setPasswordValue] = useState("");
  return (
    <Form>
      <FormGroup
        value={emailValue}
        setValue={setEmailValue}
        type="email"
        Label="Email"
        Text="We'll never share your email with anyone else."
      />
      <FormGroup
        type="password"
        Label="Password"
        Text="Password must be at least 8 characters long."
        value={passwordValue}
        setValue={setPasswordValue}
      />
      <Button variant="primary" type="submit">
        Submit
      </Button>
      <div>
        if you don't have an account, please register here:
        <Button
          variant="primary"
          onClick={() => {
            setShowRegisterForm(true);
          }}
        >
          Register
        </Button>
      </div>
    </Form>
  );
}

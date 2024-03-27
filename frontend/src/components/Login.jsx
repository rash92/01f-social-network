import FormGroup from "./FormGroup";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import {useContext, useState} from "react";
import {getJson} from "../helpers/helpers";
import AuthContext from "../store/authContext";
import Alert from "react-bootstrap/Alert";
export default function Login({setShowRegisterForm}) {
  const [emailValue, setEmailValue] = useState("");
  const [passwordValue, setPasswordValue] = useState("");
  const [error, setError] = useState({isError: false, error: ""});
  const {OnLogin} = useContext(AuthContext);
  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const res = await getJson("login", {
        method: "POST",
        credentials: "include",
        // mode: "no-cors",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({email: emailValue, password: passwordValue}),
      });

      if (!res.success) {
        setError({isError: true, error: `something went wrong ${res.status}`});
        return;
      }
      
      OnLogin(res);
    } catch (err) {
      setError({isError: true, error: err.message});
    }
  };

  return (
    <Form onSubmit={handleSubmit}>
      {error.isError && <Alert variant="danger">{error.error}</Alert>}

      <FormGroup
        value={emailValue}
        setValue={setEmailValue}
        type="email"
        Label="Email"
        Text="We'll never share your email with anyone else."
        required={true}
      />
      <FormGroup
        type="password"
        Label="Password"
        Text="Password must be at least 8 characters long."
        value={passwordValue}
        setValue={setPasswordValue}
        required={true}
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

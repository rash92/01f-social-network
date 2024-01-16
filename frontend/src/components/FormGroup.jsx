import Form from "react-bootstrap/Form";

const FormGroup = ({
  type,
  Label,
  placeholder,
  Text,
  value,
  setValue,
  required,
  accept,
  setValueImge,
}) => {
  const changeHandler = (e) => {
    if (type === "file") {
      setValue({value: e.target.value, file: e.target.files[0]});
    } else {
      setValue(e.target.value);
    }
  };

  return (
    <>
      <Form.Group className="mb-3" controlId="formBasicEmail">
        <Form.Label>
          {" "}
          {required && <span style={{color: "red"}}>*</span>} {Label}
        </Form.Label>
        {type === "textarea" ? (
          <Form.Control
            as="textarea"
            rows={3}
            placeholder={placeholder}
            value={value}
            onChange={changeHandler}
          />
        ) : (
          <Form.Control
            type={type}
            required={required}
            placeholder={placeholder}
            value={type === "file" ? value.value : value}
            onChange={changeHandler}
            accept={accept}
          />
        )}
        {/* <Form.Text className="text-muted">{Text}</Form.Text> */}
      </Form.Group>
    </>
  );
};

export default FormGroup;

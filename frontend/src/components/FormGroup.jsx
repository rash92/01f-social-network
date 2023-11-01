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
}) => {
  return (
    <>
      <Form.Group className="mb-3" controlId="formBasicEmail">
        <Form.Label>{Label}</Form.Label>
        {type === "textarea" ? (
          <Form.Control
            as="textarea"
            rows={3}
            placeholder={placeholder}
            value={value}
            onChange={(e) => setValue(e.target.value)}
          />
        ) : (
          <Form.Control
            type={type}
            required={required}
            placeholder={placeholder}
            value={value}
            onChange={(e) => setValue(e.target.value)}
            accept={accept}
          />
        )}
        <Form.Text className="text-muted">{Text}</Form.Text>
      </Form.Group>
    </>
  );
};

export default FormGroup;

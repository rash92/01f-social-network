import {Button} from "react-bootstrap";

const Action = ({numberAction, actionName, toggle}) => {
  return (
    <li className="action">
      <span>{numberAction}</span>
      <Button variant="primary" onClick={}>{actionName}</Button>
    </li>
  );
};

export default Action;
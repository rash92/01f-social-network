import {Button} from "react-bootstrap";

const Action = ({numberAction, actionName}) => {
  return (
    <li className="action">
      <span>{numberAction}</span>
      <Button variant="primary" on>{actionName}</Button>
    </li>
  );
};

export default Action;
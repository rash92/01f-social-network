import {Button} from "react-bootstrap";

const Action = ({numberAction, actionName, toggleAction}) => {
  return (
    <li className="action">
      <span>{numberAction}</span>
      <Button variant="primary" onClick={ toggleAction.n}>
        {actionName}
      </Button>
    </li>
  );
};

export default Action;
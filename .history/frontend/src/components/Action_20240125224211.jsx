import {Button} from "react-bootstrap";
import classes from "./Active.module.css";

const Action = ({numberAction, actionName, toggleAction, active}) => {
  return (
    <li className={`${classes.active} ${classes.action}`}>
      <span>{numberAction}</span>
      < button onClick={toggleAction.bind(actionName)}>{actionName}<b/>
    </li>
  );
};

export default Action;
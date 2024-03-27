import classes from "./Action.module.css";

const Action = ({numberAction, actionName, toggleAction, active}) => {
  const activeClass = actionName === active ? classes.active : "";

  return (
    <li className={`${classes.action} ${activeClass} `}>
      <span> {numberAction} </span>
      <button onClick={toggleAction.bind(null, actionName)}>
        {actionName}
      </button>
    </li>
  );
};

export default Action;

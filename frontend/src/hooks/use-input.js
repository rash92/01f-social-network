import {useState} from "react";

const useInput = (ValidateValue) => {
  const [enteredValue, setEnteredValue] = useState("");
  const [istouch, setIsTouch] = useState(false);

  const ValueIsValid = ValidateValue(enteredValue);
  let hassError = !ValueIsValid && istouch;

  const valueChangeHandler = (e) => {
    setEnteredValue(e.target.value);
  };

  const valueInputBlurHandler = (e) => {
    setIsTouch(true);
    
  };

  const reset = () => {
    setEnteredValue("");
    setIsTouch(false);
    // setIsInDB(false)
  };

  return {
    value: enteredValue,
    hassError,
    isValid: ValueIsValid,
    valueChangeHandler,
    valueInputBlurHandler,
    reset,
    // isInDataBase:isInDB,
  };
};

export default useInput;

import React, {useEffect, useContext} from "react";
import {ToastContainer, toast} from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import {Link} from "react-router-dom";
import AuthContext from "../store/authContext";

const formatUrl = (data) => {
  const options = {
    "notification requestToFollow": `/profile/${data?.ReceiverId}`,
    "notification answerRequestToFollow": `/profile/${data?.SenderId}`,
  };
  return options[data?.type] || "/";
};

const NotificationComponent = () => {
  console.log("NotificationComponent");
  const {dashBoardData} = useContext(AuthContext);
  const message =
    dashBoardData.notifications.length > 0
      ? dashBoardData.notifications[0]?.Body.Message
      : "";

  const formatedUrl = formatUrl(dashBoardData.notifications[0]);
  const id = dashBoardData?.notifications[0]?.Id;
  useEffect(() => {
    const handleNotification = () => {
      if (message) {
        toast.info(
          <Link style={{textDecoration: "none"}} to={formatedUrl}>
            {message}
          </Link>,
          {
            position: "top-center",
            autoClose: 5000, // 5 seconds
          }
        );
      }
    };

    handleNotification();
  }, [message, formatedUrl, id]);

  return (
    <div>
      <ToastContainer />
    </div>
  );
};

export default NotificationComponent;

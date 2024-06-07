import React, {useEffect, useContext, useCallback} from "react";
import {ToastContainer, toast} from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import {Link} from "react-router-dom";
import AuthContext from "../store/authContext";

const NotificationComponent = React.memo(() => {
  const formatUrl = useCallback((data) => {
    const options = {
      "notification requestToFollow": `/profile/${data?.ReceiverId}`,
      "notification answerRequestToFollow": `/profile/${data?.SenderId}`,
    };
    return options[data?.type] || "/";
  }, []);

  const {dashBoardData} = useContext(AuthContext);

  useEffect(() => {
    const handleNotification = () => {
      if (
        dashBoardData?.notifications?.length > 0 &&
        dashBoardData?.notifications[0]?.payload?.Message
      ) {
        toast.info(
          <Link
            style={{textDecoration: "none"}}
            to={formatUrl(dashBoardData.notifications[0])}
          >
            {dashBoardData.notifications[0]?.payload?.Message}
          </Link>,
          {
            position: "top-center",
            autoClose: 5000, // 5 seconds
          }
        );
      }
    };

    handleNotification();
  }, [dashBoardData?.notifications, formatUrl]);

  return (
    <div>
      <ToastContainer />
    </div>
  );
});

export default NotificationComponent;

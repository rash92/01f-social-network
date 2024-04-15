import React from "react";
import User from "./User";
import {Link} from "react-router-dom";

const SearchList = ({className, searchList}) => {
  return (
    <ul className={`${className}`}>
      {searchList?.map((el) => (
        <li key={el.id}>
          <Link to={`/profile/${el.Id}`}>
            <User Nickname={el.Nickname} Avatar={el.Avatar} />
          </Link>
        </li>
      ))}
    </ul>
  );
};

export default SearchList;

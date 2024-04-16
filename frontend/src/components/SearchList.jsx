import React from "react";
import User from "./User";
import {Link} from "react-router-dom";

const SearchList = ({className, searchList, clearSearch}) => {
  return (
    <ul className={`${className}`} onClick={clearSearch}>
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

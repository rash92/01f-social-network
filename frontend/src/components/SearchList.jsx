import React from "react";
import User from "./User";
import {Link} from "react-router-dom";

const SearchList = ({className, searchList, clearSearch, type = "user"}) => {
  return (
    <ul className={`${className}`}>
      {searchList?.map((el) => (
        <li key={el.id} onClick={clearSearch.bind(null, el.Id)}>
          {type === "user" ? (
            <Link to={`/profile/${el.Id}`}>
              <User Nickname={el.Nickname} Avatar={el.Avatar} />
            </Link>
          ) : (
            <User Nickname={el.Nickname} Avatar={el.Avatar} />
          )}
        </li>
      ))}
    </ul>
  );
};

export default SearchList;

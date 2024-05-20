import React, {useState, useEffect} from "react";
import SearchList from "./SearchList";

const Search = ({onSearch, searchList, Blur, type, addChosen}) => {
  const [query, setQuery] = useState("");
  const handleChange = (event) => {
    const {value} = event.target;
    setQuery(value);
    if (type === "follower") {
      Blur(false);
    }
  };

  useEffect(() => {
    if (query) {
      onSearch(query);
    }
  }, [query]);

  // useEffect(()=>{

  // },[])
  const clearSearch = (id) => {
    console.log(id);
    if (type === "follower") {
      addChosen(id);
    }

    setQuery("");
  };

  return (
    <div className="mb-3 position-relative">
      <div className="input-group">
        <input
          type="text"
          className="form-control"
          placeholder="Search users..."
          value={query}
          onChange={handleChange}
          onBlur={Blur?.bind(null, true)}
        />
      </div>
      {query && (
        <SearchList
          className="position-absolute top-100 start-0 w-100 p-4 bg-light border"
          searchList={searchList}
          clearSearch={clearSearch}
          type={type}
        />
      )}
    </div>
  );
};

export default Search;

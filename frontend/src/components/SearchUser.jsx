import React, {useState, useEffect} from "react";
import SearchList from "./SearchList";

const Search = ({onSearch, searchList}) => {
  const [query, setQuery] = useState("");

  const handleChange = (event) => {
    const {value} = event.target;
    setQuery(value);
  };

  useEffect(() => {
    if (query) {
      onSearch(query);
    }
  }, [query]);

  const clearSearch = () => {
    setQuery("");
  }


  return (
    <div className="mb-3 position-relative">
      <div className="input-group">
        <input
          type="text"
          className="form-control"
          placeholder="Search users..."
          value={query}
          onChange={handleChange}
        />
      </div>
      {query && (
        <SearchList
          className="position-absolute top-100 start-0 w-100 p-4 bg-light border"
          searchList={searchList}
          clearSearch={clearSearch}
        />
      )}
    </div>
  );
};

export default Search;

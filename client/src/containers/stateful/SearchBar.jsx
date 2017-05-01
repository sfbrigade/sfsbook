// Stateful react component that is responsible for rendering all instances of search bars in the app.
import React, { Component, PropTypes } from 'react';
import '../dist/stylesheets/searchbar.css';

let SearchBar = () => {
    return (
        <div className="search-bar">
            <input ref="main-search" 
                   type="text" 
                   placeholder="search here for resources">
            </input>
            <input type="submit"></input>
        </div>
    );
}

export default SearchBar;
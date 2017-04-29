// Stateful react component that is responsible for rendering all instances of search bars in the app.
import React, { Component, PropTypes } from 'react';
import '../dist/stylesheets/searchbar.css';

class SearchBar extends Component {
    constructor() {
        super();
        this.state = {
            active: false
        };
    }

    render() {
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

    handleClick() {
         const timesClicked
    }
    
}

SearchBar.propTypes = {

};

export default SearchBar;
import React, { Component, PropTypes } from 'react';

class ResourceForm extends Component {
    constructor() {
        super();
        this.state = {
            totalPosts: 0
        }
        this.handlePost = this.handlePost.bind(this);
    }

    handlePost(e) {
        e.preventDefault();
        // Validate that the correct info was inputted and then post it into the database
    }

    render() {
        return (
            <div>
                <form>
                    <span id="resource-title">
                        <p>Title:</p>
                        <input type="text" name="title"></input>
                    </span>
                    <br></br>
                    <span id="resource-body">
                        <p>Description:</p>
                        <input type="text" name="description"></input>
                    </span>
                    <input type="submit" id="submit" onClick={this.handlePost}></input>
                </form>
            </div>
        )
    }
}
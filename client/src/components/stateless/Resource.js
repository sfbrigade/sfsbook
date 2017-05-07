// This component is responsible for rendering each individual resource-containing post contributed by volunteers of the project.
import React, { Component, PropTypes } from 'react';

class Resource extends Component {
    constructor() {
        super();
        this.state = {likesCount: 0},
        this.onLike = this.onLike.bind(this);
    }

    onLike() {
        let newLikesCount = this.state.likesCount + 1;
        this.setState({likesCount: newLikesCount});
    }

    render() {
        return (
            <div>
                <h1>{this.props.header}</h1>
                <p>{this.props.body}</p>
            </div>
        )
    }
}

Resource.propTypes = {
   header: PropTypes.string,
   body: PropTypes.string
}

export default Resource;
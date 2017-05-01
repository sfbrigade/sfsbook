// The code here is quite simple, it exports a function that expects a react component as the only argument. This ChildComponent is the component you want to protect.

import React, { Component, PropTypes } from 'react';
import { bindActionCreators } from 'react-redux';
import { connect } from 'react-redux';
import * as OAUTHActions from './OAUTHActions';
import OAUTH from '../../stateless/OAUTH';
import SignIn from './SignIn';

export default (ChildComponent) => {
    class OAUTHValidation extends Component {
        static propTypes = {
            hasAuthToken: PropTypes.bool.isRequired
        };

        /* The hasAuthToken prop is the control here, if it is true then the ChildComponent will be rendered, otherwise it will render SignIn. Note, this process of rendering SignIn is fine if you don't care about SEO but you may want to redirect the user to a sign-in route if you do care about search engines indexing your protected routes. */

        render() {
            const { hasAuthToken } = this.props;
            return (hasAuthToken 
            ? <ChildComponent {...this.props} />
            : <SignIn />
            )
        }
    }
}

/* Finally OAUTHValidation is connected to the redux store session state but this could easily be switched out to use any other flux library of your choosing. Simply put, it is subscribing to changes on session. If the hasAuthToken changes in value, the component will be re-rendered if it is presently mounted. */

const mapStateToProps = ({session}) => (session);
const mapDispatchToProps(dispatch) {
    return bindActionCreators({
        ...ProductActions
    }, dispatch);
    })
}

return connect(mapStateToProps, mapDispatchToProps)(OAUTHValidation);
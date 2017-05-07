import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import * as OAUTHActions from '../../actions/OAUTHActions';

export default Auth({ user, token }) {
    // TODO: Do something with token after the h1 is rendered.
    return (
        <div className="user">
            <h1>{ user }</h1>
        </div>
    )
}
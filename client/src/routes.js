import { Route } from 'react-router';
import App from './components/containers/stateful/App';
import Login from './components/containers/stateful/Login';
import OAUTHValdiation from './components/OAUTHValdiation';

export default (
    <Route path='/' component={App}>
        <Route path='login' component={Login}>
            <Route path='oauth' component={OAUTHValdiation} />
        </Route>
    </Route>

)
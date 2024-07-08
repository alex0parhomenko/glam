import React from 'react';
import {BrowserRouter as Router, Route, Switch, Link} from 'react-router-dom';
import Profile from './Profile';
import Posts from './Posts';
import Notifications from './Notifications';

function App() {
    return (
        <Router>
            <div>
                <nav>
                    <ul>
                        <li><Link to="/profile/1">Profile</Link></li>
                        <li><Link to="/posts/1">Posts</Link></li>
                        <li><Link to="/notifications/1">Notifications</Link></li>
                    </ul>
                </nav>
                <Switch>
                    <Route path="/profile/:id" component={Profile}/>
                    <Route path="/posts/:id" component={Posts}/>
                    <Route path="/notifications/:user_id" component={Notifications}/>
                </Switch>
            </div>
        </Router>
    );
}

export default App;

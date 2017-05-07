// Stateless react component that is responsible for rendering the app's main navbar.
import React, { PropTypes } from 'react';

let Navbar = () => {
    return (
        <div>
            <nav>
                <ul>
                    <li>
                        <a>Resources</a>
                    </li>
                    <li>
                        <a>Shelter</a>
                    </li>
                    <li>
                        <a>Communities</a>
                    </li>
                    <li>
                        <a>Support Center</a>
                    </li>
                </ul>
            </nav>
        </div>
    );
}

export default Navbar;
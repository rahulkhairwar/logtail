import React from 'react';


class Header extends React.Component  {
    constructor() {
        super()
        this.state= {
            inputs: 0
        }
    }

    render() {
        return (
            <div className="header">
                <hr/>
                <h1>Logtail Corporation</h1>
            </div>
            
        )
    }
    
};


export default Header;

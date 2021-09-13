import React, {Component} from "react";
import {url1} from "./consts.js";
import './hue.css';


// LogLine

class LogLine extends Component {
    constructor() {
        super()
        this.state = {
            loading: false,
            character: []
        }
    }
    
    getUrl = (id) => {
        return url1 + id + "/"
    }

    componentDidMount() {
        this.setState({loading: true})
        
        for (let i = 1; i < 10; i++) {



            setTimeout(
                () => fetch(this.getUrl(i))
                .then(response => response.json())
                .then(data => {
                    data['number'] = i + 1;
                    this.setState(prevState => ({
                        loading: false,
                        character: [...prevState.character, data]
                    }))
                })
            , 2000)
            
        }
        
        console.log(this.state.character)
    }
    
    render() {
        // this.state.loading ? "Loading the page, please wait" :
        const characters =  this.state.character
        // const listItems = characters.map((d) => <li key={d.name}>{d.name}</li>);
        const listItems = characters.map((d) => <p key={d.name} className='hue'>{d.number} = {JSON.stringify(d, null, 2)}</p>);
        
        return (
            <div>
               {listItems}
            </div>
        )
    }
}

export default LogLine

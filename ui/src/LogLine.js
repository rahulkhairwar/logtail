import React, {Component} from "react";
import {url1} from "./consts.js";



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



           
                fetch(this.getUrl(i))
                .then(response => response.json())
                .then(data => {
                    data['number'] = i + 1;
                    this.setState(prevState => ({
                        loading: false,
                        character: [...prevState.character, data]
                    }))
                })
            
            
        }
        
        
        console.log(this.state.character)
    }
   
    listItems = () => {
        const characters =  this.state.character
        let timer = 0;
        for (let i = 0; i < characters.length; i++) {
          setTimeout(() => document.getElementById('data').innerHTML = characters[i].name, timer);
          timer = timer + 1000;
        }
      }
    
    render() {
        // this.state.loading ? "Loading the page, please wait" :
        
        // const listItems = characters.map((d) => <li key={d.name}>{d.name}</li>);
                
        this.listItems()
        
        
        return (
            <div className="Logline">
                <div >
                    <input type="text" placeholder="Search for the data here" className="input"></input>
                    <button className="button">Search</button>
                </div>
                
                <div id="data" ></div>
            </div>
        )
    }
}

export default LogLine


import React from 'react';
import ReactDOM from 'react-dom';
import "!style!css!sass!./scss/main.scss";
import update from 'react-addons-update'
import {Router, Route, IndexRoute, Link, browserHistory} from 'react-router'
import 'whatwg-fetch';

import StockList from './myStock.js';
import UserStockList from './userStockList.js';
import Main from './main.js';

const API_DOMAIN = "http://localhost:8080/";
const API_SEARCH = "search";
const API_HEADER = {
	'Accept': 'application/json',
	'Content-Type': 'application/json'
};

var Home = React.createClass({
	componentDidMount: function(){
		this.nav = document.getElementsByClassName("tab");
	},

	clickedOnTab: function(e){
		Array.from(this.nav).forEach(function(target){
			target.className = "tab";
		});

		e.target.className="tab active";
	},

	render: function () {
		let Footer = function(){
			return(
				<footer>
					© Created 2016 By Yen-Chieh Chen.
				</footer>
			)
		};

		return (
			<div>
				<div id="nav">
					<ul className="nav nav-tabs">
						<li role="presentation" className="tab active" onClick={this.clickedOnTab}><Link to="/">Search</Link></li>
						<li role="presentation" className="tab" onClick={this.clickedOnTab}><Link to="/myStock">My Stock</Link></li>
					</ul>
				</div>

				{this.props.children}

				<Footer/>
			</div>
		)
	}
});

ReactDOM.render(
	<Router history={browserHistory}>
		<Route path="/" component={Home}>
			<IndexRoute component={Main}/>
			<Route path="myStock" component={UserStockList} />
		</Route>
	</Router>
	, document.getElementById('app'));
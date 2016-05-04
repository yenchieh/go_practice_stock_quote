
import React from 'react';
import ReactDOM from 'react-dom';
import "!style!css!sass!./scss/main.scss";
import update from 'react-addons-update'
import {Router, Route, IndexRoute, Link, browserHistory} from 'react-router'
import 'whatwg-fetch';

import UserStockList from './userStockList.js';
import QuoteResultList from './quoteResultList.js';
import Config from './mainConfig.js'

var Home = React.createClass({

	getInitialState: function() {
		var storage = window.localStorage.getItem("Main") || "{}";
		if(storage) {
			storage = JSON.parse(storage);

		}

		//Check if user exist
		if(storage.userName && storage.userEmail){
			this.checkUser(storage.userName, storage.userEmail);
		}


		return {
			userName: storage.userName,
			userEmail:storage.userEmail
		}
	},

	checkUser: function(name, email){
		fetch(Config.API_DOMAIN + Config.API_CHECK_USER + "?userName=" + name + "&userEmail=" + email, {
			method: 'GET',
			headers: Config.API_JSON_HEADER
		}).then((response) => {
			if(response.status == 404){

				this.setState({
					userName: "",
					userEmail: ""
				});
			}

		});
	},

	componentDidMount: function(){
		this.nav = document.getElementsByClassName("tab");
	},

	clickedOnTab: function(e){
		Array.from(this.nav).forEach(function(target){
			target.className = "tab";
		});

		e.target.className="tab active";
	},

	updateUserInfo: function(data){
console.error(data);
	},

	render: function () {
		let Footer = function(){
			return(
				<footer>
					© Created 2016 By Yen-Chieh Chen.
				</footer>
			)
		};

		let customStockLink = {
			pathname: "/userStockList",
			query: {
				userName: this.state.userName,
				userEmail: this.state.userEmail
			}
		};

		let renderTab = function(){
			return (
				<li role="presentation" className="tab" onClick={this.clickedOnTab}><Link to={customStockLink}>{this.state.userName}</Link></li>
			)
		}.bind(this);

		return (
			<div>
				<div id="nav">
					<ul className="nav nav-tabs">
						<li role="presentation" className="tab active" onClick={this.clickedOnTab}><Link to="/">Search</Link></li>
						{renderTab()}
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
			<IndexRoute component={QuoteResultList}/>
			<Route path="userStockList" component={UserStockList} />
		</Route>
	</Router>
	, document.getElementById('app'));
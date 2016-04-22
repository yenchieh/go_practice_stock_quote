/**
 * Created by yen-chieh on 4/20/16.
 */
import React from 'react';
import ReactDOM from 'react-dom';
import QuoteTable from './component/quote.js';
import update from 'react-addons-update'
import {Router, Route, Link, browserHistory} from 'react-router'
import 'whatwg-fetch';

const API_DOMAIN = "http://localhost:8080/";
const API_REMOVE = "removeFromList";
const API_USER_LIST = "getUserStockList";
const API_HEADER = {
	'Accept': 'application/json',
	'Content-Type': 'application/json'
};

var UserStockList = React.createClass({
	getInitialState: function () {
		return {
			title: "",
			quoteData: []
		}
	},

	componentDidMount: function(){
		fetch(API_DOMAIN + API_USER_LIST + "?userName=Jay", {
			method: 'get',
			header: API_HEADER
		}).then((response) => response.json())
		.then((data) => {
				this.fetchUserList(data.query);
			}
		);
	},

	fetchUserList: function(data){
		console.error(data);
		if(data.count == 0){
			this.setState({title: "You have no data"});
			return;
		}

		data.results.quote.map((q, i) =>{
			var updatedQ = update(q, {$merge: {buttonType: "remove"}});
			data.results.quote = update(data.results.quote, {[i]: {$set: updatedQ}})
		});

		this.setState({quoteData: data.results.quote});
	},

	removeFromStockList: function(key){
		var username = "Jay";

		var data = {
			listId: this.state.quoteData[key].id,
			username: username

		};
		console.error(data);
		fetch(API_DOMAIN + API_REMOVE, {
			method: 'POST',
			body: JSON.stringify(data),
			headers: API_HEADER
		}).then((response) => {
			if(response.ok){
				this.removeFromStockList_load(key);
			}
		});
	},

	removeFromStockList_load: function(key){
		let quoteData = this.state.quoteData;

		quoteData = update(quoteData, {$splice: [[key, 1]]});

		this.setState({quoteData: quoteData});
	},

	render: function(){
		return (
			<div>
				<h1>{this.state.title}</h1>
				<QuoteTable data={this.state.quoteData} removeButtonCallback={this.removeFromStockList}/>
			</div>

		)
	}
});


export default UserStockList;
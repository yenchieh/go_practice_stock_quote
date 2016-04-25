/**
 * Created by yen-chieh on 4/20/16.
 */
import React from 'react';
import ReactDOM from 'react-dom';
import QuoteTable from './component/quote.js';
import update from 'react-addons-update';
import {Router, Route, Link, browserHistory} from 'react-router';
import LocalStorageMixin from 'react-localstorage';
import 'whatwg-fetch';

import Config from './mainConfig.js';

var UserStockList = React.createClass({
	mixins: [LocalStorageMixin],

	getInitialState: function () {
		return {
			title: "",
			quoteData: [],
			userName: "",
			userEmail: ""
		}
	},

	componentDidMount: function(){
		let query = this.props.location.query;
		this.setState({
			userName: query.userName,
			userEmail: query.userEmail
		});
		fetch(Config.API_DOMAIN + Config.API_USER_STOCK_LIST + "?userName="+query.userName+"&userEmail="+query.userEmail, {
			method: 'get',
			header: Config.API_JSON_HEADER
		}).then((response) => response.json())
		.then((data) => {
				this.fetchUserList(data.query);
			}
		);
	},

	fetchUserList: function(data){
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

		fetch(Config.API_DOMAIN + Config.API_REMOVE_LIST, {
			method: 'POST',
			body: JSON.stringify(data),
			headers: Config.API_JSON_HEADER
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
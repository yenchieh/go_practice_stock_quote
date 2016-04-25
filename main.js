import React from 'react';
import ReactDOM from 'react-dom';
import "!style!css!sass!./scss/main.scss";
import QuoteTable from './component/quote.js';
import StockList from './myStock.js';
import Dialog from './tools/dialog.js';
import update from 'react-addons-update';
import {Router, Route, Link, browserHistory} from 'react-router';
import LocalStorageMixin from 'react-localstorage';
import 'whatwg-fetch';

import Config from './mainConfig.js'

var Main = React.createClass({
	mixins: [LocalStorageMixin],
	getInitialState: function () {
		return {
			quoteData: [],
			pop: {},
			showDialog: false,
			userName: "",
			userEmail: ""
		}
	},
	getStateFilterKeys: function() {
		return ['quoteData', 'pop', 'userName', 'userEmail'];
	},
	componentDidMount: function () {
		this.symbolInput = document.getElementById("symbolSearchInput");
	},
	searchSymbol: function () {
		fetch(Config.API_DOMAIN + Config.API_SEARCH_STOCK + "?symbol=" + this.symbolInput.value, {
			method: 'GET',
			headers: Config.API_JSON_HEADER
		}).then((response) => response.json())
			.then((data) => {
				if (!data || data.length == 0) {
					return false;
				}
				this.symbolInput.value = "";
				var quoteData = this.state.quoteData;
				let newQuote = data.query.results.quote;
				newQuote = update(newQuote, {$merge: {buttonType: "add"}});
				quoteData = update(quoteData, {$unshift: [newQuote]});
				this.setState({quoteData: quoteData});
			}
		);

	},
	keypressed: function (e) {
		if (e.key == "Enter") {
			this.searchSymbol();
		}
	},

	clickRemoveButton: function (key) {
		let quoteData = this.state.quoteData;
		quoteData = update(quoteData, {$splice: [[key, 1]]});
		this.setState({quoteData: quoteData});
	},

	clickAddButton: function (quote, event) {
		if (!this.state.userName || !this.state.userEmail) {
			this.showDialog();
			return;
		}

		var clickedButton = $(event.target);

		if (clickedButton.hasClass("disabled")) {
			return;
		}

		clickedButton.addClass("disabled");
		this.storeStock({
			button: clickedButton,
			quote: quote
		});
	},

	storeStock: function (oArgs) {
		oArgs = oArgs || {};
		var clickedButton = oArgs.button;
		var userName = this.state.userName;
		var userEmail = this.state.userEmail
		var data = {
			symbol: oArgs.quote.symbol,
			stockName: oArgs.quote.name,
			userName: userName,
			userEmail: userEmail
		};

		fetch(Config.API_DOMAIN + Config.API_ADD_STOCK, {
			method: 'POST',
			body: JSON.stringify(data),
			headers: Config.API_JSON_HEADER
		}).then((response) => response.json())
			.then((data) => {
				if (!data || data.length == 0) {
					clickedButton.removeClass("disabled");

				} else {
					this.stockAdded(clickedButton);
				}
			}
		);
	},

	stockAdded: function (clickedButton) {
		clickedButton.parents("tr").css({
			"background-color": "gainsboro"
		});

		let index = clickedButton.attr("data-id");
		let updateButton = update(this.state.quoteData, {
			[index]: {
				buttonType: {
					$set: "remove"
				}
			}
		});

		this.setState({
			quoteData: updateButton
		});

	},

	showDialog: function () {
		this.setState({
			showDialog: true,
			pop: {
				headerText: "Watch List Portfolio",
				bodyText: "Create watch list portfolio by enter information below",
				confirmButton: "Register / Login",
				additionalBody: this.accountDialog()
			}
		})
	},

	closedDialog: function () {
		this.setState({showDialog: false});
	},

	confirmCallback: function($ele) {
		var name = $ele.find('#userNameInput').val();
		var email = $ele.find('#userEmailInput').val();

		this.setState({
			userName: name,
			userEmail: email
		});

	},

	accountDialog: () => {
		return (
			<div className="accountDialog">
				<input type="text" name="name" className="form-control" id="userNameInput" placeholder="Name"/><br/>
				<input type="text" name="email" className="form-control" id="userEmailInput" placeholder="Email"/>
			</div>
		)

	},


	render: function () {
		let showDialog = () => {
			if (this.state.showDialog) {
				return (
					<Dialog headerText={this.state.pop.headerText} bodyText={this.state.pop.bodyText}
									additionalBody={this.state.pop.additionalBody} confirmButton={this.state.pop.confirmButton}
									closedCallback={this.closedDialog} confirmCallback={this.confirmCallback}/>
				)
			}
		};

		return (
			<div>
				<div id="mainSearchComponent">
					<figure className="highlight">
						<div className="form-inline">
							<div className="form-group">
								<label for="symbolSearchInput" className="sr-only">Enter Stock Here</label>
								<input type="text" name="symbol" className="form-control" id="symbolSearchInput"
											 placeholder="Enter symbol or stock name" onKeyPress={this.keypressed}/>
							</div>
							<button className="btn btn-sm btn-primary" onClick={this.searchSymbol}>Search</button>
						</div>
					</figure>
				</div>

				<div id="quoteList">
					<h1>Quote List</h1>
					<QuoteTable data={this.state.quoteData} removeButtonCallback={this.clickRemoveButton}
											addButtonCallback={this.clickAddButton}/>
				</div>

				<StockList/>
				{showDialog()}
			</div>

		)
	}
});

export default Main
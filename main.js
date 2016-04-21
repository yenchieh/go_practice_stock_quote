
var React = require('react');
var ReactDOM = require('react-dom');
var css = require("!style!css!sass!./scss/main.scss");
var QuoteTable = require('./component/quote.js');
var StockList = require('./myStock.js');
var update = require('react-addons-update');
import 'whatwg-fetch';

const API_DOMAIN = "http://localhost:8080/";
const API_SEARCH = "search";
const API_HEADER = {
	'Accept': 'application/json',
	'Content-Type': 'application/json'
};

var Main = React.createClass({
	getInitialState: function () {
		return {
			domain: 'http://localhost:8080/',
			searchPath: 'search',
			quoteData: [],
			searchIndex: []
		}
	},

	componentDidMount: function () {
		this.symbolInput = document.getElementById("symbolSearchInput");
	},

	searchSymbol: function(){
		fetch(API_DOMAIN + API_SEARCH + "?symbol=" + this.symbolInput.value, {
			method: 'GET',
			mode: 'cors',
			headers: API_HEADER
		}).then((response) => response.json())
			.then((data) => {

				if(!data || data.length == 0){
					return false;
				}
				this.symbolInput.value = "";
				var quoteData = this.state.quoteData;
				let newQuote = data.query.results.quote;
				newQuote = update(newQuote, {$merge: {buttonType: "add"}});
				quoteData.push(newQuote);
				this.setState({quoteData: quoteData});
			}
		);

	},

	keypressed: function(e){

		if(e.key == "Enter"){
			this.searchSymbol();
		}
	},

	clickRemoveButton: function(key){
		console.error(key);
	},

	render: function () {
		return (
			<div>
				<div id="nav">
					<ul className="nav nav-tabs">
						<li role="presentation" className="active"><a href="#">Search</a></li>
						<li role="presentation"><a href="#">My Stock</a></li>
					</ul>
				</div>
				<div id="mainSearchComponent">
					<figure className="highlight">
						<div className="form-inline">
							<div className="form-group">
								<label for="symbolSearchInput" className="sr-only">Enter Stock Here</label>
								<input type="text" name="symbol" className="form-control" id="symbolSearchInput"
											 placeholder="Enter symbol or stock name"  onKeyPress={this.keypressed}/>
							</div>
							<button className="btn btn-sm btn-primary" onClick={this.searchSymbol}>Search</button>
						</div>
					</figure>
				</div>

				<div id="quoteList">
					<h1>Quote List</h1>
					<QuoteTable data={this.state.quoteData} removeButtonCallback={this.clickRemoveButton}/>
				</div>

				<StockList/>
			</div>

		)
	}
});

ReactDOM.render(<Main/>, document.getElementById('app'));
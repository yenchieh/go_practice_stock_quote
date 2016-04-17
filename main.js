
var React = require('react');
var ReactDOM = require('react-dom');
var css = require("!style!css!sass!./scss/main.scss");

var QuoteTable = require('./component/quote.js');


var Main = React.createClass({
	getInitialState: function () {
		return {
			domain: 'http://localhost:8080/',
			searchPath: 'search',
			$app: $('div#app'),
			quoteData: [],
			searchIndex: []
		}
	},

	componentDidMount: function () {
		this.$symbolInput = $('#symbolSearchInput');
	},

	searchSymbol: function(){
		var symbolInput = $('input#symbolSearchInput', this.state.$app);
		if(!symbolInput){
			return false;
		}

		$.ajax({
			url: this.state.domain + this.state.searchPath,
			data: {
				symbol: symbolInput.val()
			},
			dataType: 'json',
			success: function(data){
				if(!data || data.length == 0){
					return false;
				}
				this.$symbolInput.val("");
				var quoteData = this.state.quoteData;
				quoteData.push(data.query.results.quote);
				this.setState({quoteData: quoteData});
			}.bind(this)
		});
	},

	keypressed: function(e){

		if(e.key == "Enter"){
			this.searchSymbol();
		}
	},

	render: function () {
		return (
			<div>
				<div id="nav">
					<nav className="navbar navbar-default navbar-static">
						<a className="navbar-toggle" data-toggle="collapse" data-target=".navbar-collapse">
							<span className="glyphicon glyphicon-bar">aa</span>
							<span className="glyphicon glyphicon-bar"></span>
							<span className="glyphicon glyphicon-bar"></span>
						</a>
					</nav>
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
					<QuoteTable data={this.state.quoteData}/>
				</div>
			</div>

		)
	}
});

ReactDOM.render(<Main/>, document.getElementById('app'));
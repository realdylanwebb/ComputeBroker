import { AppBar, Box, Card, CardActions, CardContent, Container, Button, TextField, Typography, Divider, Toolbar, makeStyles, Grid, Paper } from '@material-ui/core';
import React from 'react';

class App extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      apiKey: "hlka",
      email: "dylan@test.com",
      register: true
    };

    this.displayLogin = this.displayLogin.bind(this);
    this.displayRegister = this.displayRegister.bind(this);
    this.updateKey = this.updateKey.bind(this);
    this.updateEmail = this.updateEmail.bind(this);
  }

  displayLogin() {
    this.setState({register: false})
  }

  displayRegister() {
    this.setState({register: true})
  }

  updateKey(key) {
    this.setState({apiKey: key})
  }

  updateEmail(email) {
    this.setState({email: email})
  }

  render() {
    let status;
    if (!this.state.apiKey) {
      if (this.state.register) {
        status = (
          <Register login={this.displayLogin} updateKey={this.updateKey} updateEmail={this.updateEmail} />
        )
      } else {
        status = (
          <Login register={this.displayRegister} updateKey={this.updateKey} updateEmail={this.updateEmail} />
        )
      } 
    } else {
      status = (
        <Dashboard apiKey={this.state.apiKey}/>
      )
    }
    
    return (
      <React.Fragment>
        <Navbar email={this.state.email} updateKey={this.updateKey} updateEmail={this.updateEmail} />
        {status}
      </React.Fragment>
    )

  }
}

class Navbar extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {

    let button
    if (this.props.email != null) {
      button = (
        <React.Fragment>
          <Button color="inherit" onClick={()=>{this.props.updateEmail(null); this.props.updateKey(null);}}>Logout</Button>
        </React.Fragment>
      )
    }

    return (
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h5">ComputeBroker</Typography>
          {button}
        </Toolbar>
      </AppBar>
    )
  }
}


class Register extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      email : "",
      password : "",
      hostAddress: "",
      publicKey: ""
    };

    this.onFormUpdate = this.onFormUpdate.bind(this)
    this.onFormSubmit = this.onFormSubmit.bind(this)
  }

  onFormUpdate(event) {
    switch(event.target.name) {
      case "emailInput":
        this.setState({email: event.target.value});
        break;
      case "passwordInput":
        this.setState({password: event.target.value});
        break;
      case "addressInput":
        this.setState({hostAddress: event.target.value});
        break;
      case "publicKeyInput":
        this.setState({publicKey: event.target.value});
        break;
    }
  }

  onFormSubmit() {

    let currentEmail = this.state.email

    fetch("http://localhost:8080/client", {
      method: "POST",
      cache: "no-cache",
      headers: {
        "Content-Type" : "application/json"
      },
      body: JSON.stringify({
        email: this.state.email,
        password: this.state.password,
        pubKey: this.state.publicKey,
        address: this.state.hostAddress
      })
    })
    .then(data => {
      if ((200 <= data.status) && (data.status < 300)) {
        let res = data.json()
        fetch("http://localhost:8080/login", {
          method: "POST",
          cache: "no-cache",
          headers: {
            "Content-Type" : "application/json"
          },
          body: JSON.stringify({
            email: res.email,
            password: res.password
          })
        })
        .then(data => {
          if ((200 <= data.status) && (data.status < 300)) {
            let res = data.json()
            this.props.updateKey(res.token);
            this.props.updateEmail(currentEmail);
          }
        })
        .catch(err => {
          console.error(err)
        });
      }
    })
    .catch(err => {
      console.error(err)
    })
  }

  render() {
    return (
      <Box my={8} >
        <Container maxWidth="sm">
          <Card>
            <CardContent>
              <Typography variant="h3">Register</Typography>
              <TextField label="Email" name="emailInput" required type="email" fullWidth="true" onChange={this.onFormUpdate}/>
              <TextField label="Password" name="passwordInput" required type="password" fullWidth="true" onChange={this.onFormUpdate}/>
              <TextField label="Host Address" name="addressInput" required fullWidth="true" onChange={this.onFormUpdate}/>
              <TextField label="PublicKey" name="publicKeyInput" required fullWidth="true" onChange={this.onFormUpdate}/>
            </CardContent>
            <CardActions>
              <Button size="small" color="primary" onClick={this.onFormSubmit}>Submit</Button>
            </CardActions>
            <Divider/>
            <CardActions>
              <Typography>Already have an account?</Typography>
              <Button size="small" color="primary" onClick={()=>{this.props.login()}}>Log in</Button>
            </CardActions>
          </Card>
        </Container>
      </Box>
    )
  }
}


class Login extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      email: "",
      password: "",
    };

    this.onFormUpdate = this.onFormUpdate.bind(this);
    this.onFormSubmit = this.onFormSubmit.bind(this);
  }

  onFormUpdate(event) {
    switch(event.target.name) {
      case "emailInput":
        this.setState({email: event.target.value});
        break;
      case "passwordInput":
        this.setState({password: event.target.value});
        break;
    }
  }

  onFormSubmit(event) {

    let currentEmail = this.state.email

    fetch("http://localhost:8080/login", {
      method: "POST",
      cache: "no-cache",
      headers: {
        "Content-Type" : "application/json"
      },
      body: JSON.stringify({email: this.state.email, password: this.state.password})
    })
    .then(data => {
      if ((200 <= data.status) && (data.status < 300)) {
        let res = data.json()
        this.props.updateKey(res.token);
        this.props.updateEmail(currentEmail);
      }
    })
    .catch(err => {
      console.error(err)
    })
  
  }

  render() {
    return (
      <Box my={8}>
        <Container maxWidth="sm">
          <Card>
            <CardContent>
              <Typography variant="h3">Login</Typography>
              <TextField label="Email" required type="email" name="emailInput" fullWidth="true" onChange={this.onFormUpdate}/>
              <TextField label="Password" required type="password" name="passwordInput" fullWidth="true" onChange={this.onFormUpdate}/>
            </CardContent>
            <CardActions>
              <Button size="small" color="primary" onClick={this.onFormSubmit}>Submit</Button>
            </CardActions>
            <Divider/>
            <CardActions>
              <Button size="small" color="primary" onClick={()=>{this.props.register()}}>Register</Button>
            </CardActions>
          </Card>
        </Container>
      </Box>
    )
  }
}


class Dashboard extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      email: "",
      pubKey: "",
      address: "",
      jobsAvailable: 0
    }

    this.updateJobs = this.updateJobs.bind(this)
  }

  updateJobs(value) {
    this.setState({jobsAvailable: value})
  }

  componentDidMount() {
    fetch("http://localhost:8080/client/get", {
      method: "POST",
      cache: "no-cache",
      headers: {
        "Content-Type" : "application/json",
        "Authorization" : this.props.apiKey
      },
    })
    .then(data => {
      if ((200 <= data.status) && (data.status < 300)) {
        let res = data.json()
        this.setState({email: res.email, pubKey: res.pubKey, address: res.address, jobsAvailable: res.jobsAvailable})
      }
    })
    .catch(err => {
      console.error(err)
    })
  }

  render() {
    return (
      <Box my={2}>
        <Grid container spacing={3}>
          <Grid item xs={6}>
            <SignalReadyAction jobsAvailable={this.state.jobsAvailable} updateJobs={this.updateJobs} apiKey={this.props.apiKey}/>
          </Grid>
          <Grid item xs={6}>
            <WorkerGroups apiKey={this.props.apiKey}/>
          </Grid>
        </Grid>
      </Box>
    )
  }

}


class SignalReadyAction extends React.Component {
  constructor(props) {
    super(props)
    this.onSubmitReady = this.onSubmitReady.bind(this)
  }

  onSubmitReady(ev) {
    let toggle = "1"
    if (this.props.jobsAvailable != 0) {
      toggle = "0"
    }

    fetch("http://localhost:8080/client/signal/"+toggle, {
      method: "POST",
      cache: "no-cache",
      headers: {
        "Content-Type" : "application/json",
        "Authorization" : this.props.apiKey
      },
    })
    .then(data => {
      if ((200 <= data.status) && (data.status < 300)) {
        this.props.updateJobs(parseInt(toggle))
      }
    })
    .catch(err => {
      console.error(err)
    })
  }

  render() {
    let text, text2
    if (this.props.jobsAvailable != 0) {
      text = "You are currently available for a workload"
      text2 = "Set Available"
    } else {
      text = "You are not available for a workload"
      text2 = "Set Unavailable"
    }
    return (
      <Card>
        <Card>
            <CardContent>
              <Typography variant="h3">{text}</Typography>
            </CardContent>
            <CardActions>
              <Button size="small" color="primary" onClick={this.onSubmitReady}>{text2}</Button>
            </CardActions>
          </Card>
      </Card>
    )
  }
}


class WorkerGroups extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      sessions: [],
      numWorkers: 0
    }
  }

  componentDidMount() {
    fetch("http://localhost:8080/client/get", {
      method: "POST",
      cache: "no-cache",
      headers: {
        "Content-Type" : "application/json",
        "Authorization" : this.props.apiKey
      },
    })
    .then(data => {
      if ((200 <= data.status) && (data.status < 300)) {
        let res = data.json()
        this.setState({sessions: res.sessions})
      }
    })
    .catch(err => {
      console.error(err)
    })
  }

  render() {
    let groups = this.state.sessions.map((session)=>WorkerGroup({workers:session.workers}))
    return(
      <Card>
        <CardContent>
         {groups} 
        </CardContent>
        <CardActions>
        </CardActions>
      </Card>
    )
  }
}


function WorkerGroup(props) {
  let elements = props.workers.map((worker)=>Worker({worker: worker}))
  return (
    <Paper>
      <Grid container>
        {elements}
      </Grid>
    </Paper>
  )
}

function Worker(props) {
  return(
    <Grid item>
      <Paper>
        <Typography>{props.worker.address}</Typography>
        <Typography>{props.worker.pubKey}</Typography>
      </Paper>
    </Grid>
  )
}


export default App;

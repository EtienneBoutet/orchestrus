import React, { useState, useEffect } from 'react';
import ReactDOM from 'react-dom';
import { makeStyles, createMuiTheme } from '@material-ui/core/styles';
import { ThemeProvider } from '@material-ui/styles';
import Typography from '@material-ui/core/Typography';
import Collapse from '@material-ui/core/Collapse';
import Box from '@material-ui/core/Box';
import Button from '@material-ui/core/Button';
import IconButton from '@material-ui/core/IconButton';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import TextField from '@material-ui/core/TextField';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import Checkbox from '@material-ui/core/Checkbox';
import axios from 'axios'
import Paper from '@material-ui/core/Paper';
import { 
  Refresh,
  KeyboardArrowUp,
  KeyboardArrowDown,
  ControlPoint,
  ArrowRightAlt,
  Delete
} from '@material-ui/icons';
import './index.css';

const theme = createMuiTheme({
  palette: {
    primary: {
      light: '#FF8C00',
      main: '#FF8C00',
      dark: '#FF8C00',
      contrastText: '#fff',
    },
    type: 'dark',
  },
  typography: {
    fontFamily: [
      'Inter',
    ].join(','),
  },
});

const useStyles = makeStyles({
  workerRow: {
    '& > *': {
      borderBottom: 'unset',
    },
  },
  list: {
    overflow: 'auto',
    position: 'relative',
    minWidth: 150,
    maxHeight: 30
  },
  workerTable: {
    minWidth: 700
  }
});

function AddWorkerDialog(props) {
  const { openAddWorker, toggleOpenAddWorker, workers, setWorkers } = props;
  const [worker, setWorker] = useState({"ip": "", "active": true, "images": []});

  function addWorker(worker) {
    axios.post("http://localhost:1235/workers", worker)
      .then(res => {
        setWorkers(workers.concat(res.data))
      })
      .catch(err => {
        alert(err.response.data)
      });
  }

  return (
    <Dialog open={openAddWorker} onClose={toggleOpenAddWorker}>
      <DialogTitle>Add Worker</DialogTitle>
      <DialogContent>
        <TextField
          variant="outlined"
          label="IP"
          name="ip"
          value={worker.ip}
          onChange={(e) => setWorker({...worker, [e.target.name]: e.target.value})}
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={() => {addWorker(worker); toggleOpenAddWorker();}} color="primary">
          Save
        </Button>
        <Button onClick={toggleOpenAddWorker} color="primary">
          Cancel
        </Button>
      </DialogActions>
    </Dialog>
  )
}

function renderPorts(ports) {
  let portsEl = []
  for (let key in ports) {
    if (ports.hasOwnProperty(key)) {
      portsEl.push(
        <ListItem key={key}>
          <div class="port-list">
            <div class="port-list-item">
              <span>{key}</span>
            </div>
            <ArrowRightAlt/>
            <div class="port-list-item">
              <span>{ports[key]}</span>
            </div>
          </div>
        </ListItem>
      );
    }
  }
  return portsEl;
}

function PortsWidget(props) {
  const { ports, setPorts } = props;

  function handlePortChange(e, i) {
    let copy = [...ports];
    copy[i] = {...ports[i], [e.target.name]: e.target.value};
    setPorts(copy);
  }

  return (
    <div>
      <div class="ports-header">
        <Typography variant="h6" gutterBottom component="div">
          Ports
        </Typography>
        <IconButton onClick={() => setPorts(ports.concat({"local": "", "public": ""}))}>
          <ControlPoint color="primary"/>
        </IconButton>
      </div>
      {ports.map((port, i) => (
        <div class="ports-row">
          <TextField
            variant="outlined"
            label="Local"
            name="local"
            type="number"
            value={port.local}
            onChange={(e) => handlePortChange(e, i)}
          />
          <ArrowRightAlt/>
          <TextField
            variant="outlined"
            label="Public"
            name="public"
            type="number"
            value={port.public}
            onChange={(e) => handlePortChange(e, i)}
          />
        </div>
      ))}
    </div>
  )
}

function AddImageDialog(props) {
  const { openAddImage, toggleOpenAddImage, selectWorker, workers, setWorkers } = props;
  const [ports, setPorts] = useState([{"local": "", "public": ""}]);
  const [image, setImage] = useState({"name": ""});

  function addImage(image) {
    axios.post("http://localhost:1235/images", image)
        .then(res => {
            console.log(res)
            let copy = [...workers]
            copy[selectWorker].images.push(res.data);
            setWorkers(copy);
        })
        .catch(err => {
            alert(err.response.data)
        })

  }

  function formatPorts(ports) {
    let formatedPorts = {}
    ports.forEach(port => {
      formatedPorts[port.local] = port.public
    });
    return formatedPorts;
  }

  return (
    <Dialog open={openAddImage} onClose={toggleOpenAddImage}>
      <DialogTitle>Add image to {workers[selectWorker].ip}</DialogTitle>
      <DialogContent>
        <TextField
          variant="outlined"
          label="Name"
          name="name"
          value={image.name}
          onChange={(e) => setImage({...image, [e.target.name]: e.target.value})}
        />
        <PortsWidget ports={ports} setPorts={setPorts}/>
      </DialogContent>
      <DialogActions>
        <Button onClick={() => {addImage({...image, "ports": formatPorts(ports)}); toggleOpenAddImage();}} color="primary">
          Save
        </Button>
        <Button onClick={toggleOpenAddImage} color="primary">
          Cancel
        </Button>
      </DialogActions>
    </Dialog>
  )
}

function WorkerTableRow(props) {
  const { worker, toggleOpenAddImage, setSelectWorker, index, workers, setWorkers } = props;
  const [open, setOpen] = useState(false);
  const classes = useStyles();

  function handleRowClick() {
    setOpen(!open)
    setSelectWorker(index)
  }

  function deleteImage(j) {
    axios.delete("http://localhost:1235/workers/" + worker.ip + "/images/" + worker.images[j].id)
        .catch(err => {
            alert(err.response.data)
        })
        .then(() => {
            let copy = [...workers];
            delete copy[index].images[j];
            setWorkers(copy);
        })
  }

  return (
    <React.Fragment>
      <TableRow className={classes.workerRow}>
        <TableCell>
          <IconButton aria-label="expand row" size="small" onClick={handleRowClick}>
            {open ? <KeyboardArrowUp /> : <KeyboardArrowDown />}
          </IconButton>
        </TableCell>
        <TableCell component="th" scope="row">
          {worker.ip}
        </TableCell>
        <TableCell>
          <Checkbox
            checked={worker.active}
            name="active"
            color="primary"
            readOnly
          />
        </TableCell>
      </TableRow>
      <TableRow>
        <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
          <Collapse in={open} timeout="auto" unmountOnExit>
            <Box margin={1}>
              <div class="images-header">
                <Typography variant="h6" gutterBottom component="div">
                  Images
                </Typography>
                <IconButton onClick={toggleOpenAddImage}>
                  <ControlPoint color="primary"/>
                </IconButton>
              </div>
              <Table size="small">
                <TableHead>
                  <TableRow>
                    <TableCell>Name</TableCell>
                    <TableCell>ID</TableCell>
                    <TableCell>Ports</TableCell>
                    <TableCell/>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {worker.images.map((image, j) => (
                    <TableRow key={image.id}>
                      <TableCell component="th" scope="row">
                        {image.name}
                      </TableCell>
                      <TableCell>{image.id}</TableCell>
                      <TableCell>
                        <List className={classes.list}>
                          {renderPorts(image.ports)}
                        </List>
                      </TableCell>
                      <TableCell>
                        <IconButton onClick={() => deleteImage(j)}>
                          <Delete/>
                        </IconButton>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </Box>
          </Collapse>
        </TableCell>
      </TableRow>
    </React.Fragment>
  );
}

function App() {
  const [openAddImage, setOpenAddImage] = useState(false);
  const [openAddWorker, setOpenAddWorker] = useState(false);
  const [selectWorker, setSelectWorker] = useState(null);
  const [workers, setWorkers] = useState([]);
  const classes = useStyles();
  
  function refreshWorkers() {
    axios.get("http://localhost:1235/workers")
        .then(resp => {
            setWorkers(resp.data)
        })    
        .catch(err => {
            setWorkers([])
        })

  }

  function toggleOpenAddImage(){
    setOpenAddImage(!openAddImage)
  }

  function toggleOpenAddWorker(){
    setOpenAddWorker(!openAddWorker)
  }

  useEffect(() => {
    refreshWorkers();
  }, []);

  return (
    <ThemeProvider theme={theme}>
      {selectWorker !== null ?
      <AddImageDialog
        openAddImage={openAddImage}
        toggleOpenAddImage={toggleOpenAddImage}
        selectWorker={selectWorker}
        setWorkers={setWorkers}
        workers={workers}
      /> : ""}
      <AddWorkerDialog
        openAddWorker={openAddWorker}
        toggleOpenAddWorker={toggleOpenAddWorker}
        workers={workers}
        setWorkers={setWorkers}
      />
      <Paper elevation={3}>
        <div class="container">
          <div class="images-header">
            <Typography variant="h2" component="h1" color="primary">
              Orchestrus
            </Typography>
            <IconButton onClick={refreshWorkers}>
              <Refresh/>
            </IconButton>
          </div>
          <br/>
          <div class="images-header">
            <Typography variant="h4" component="h1">
              Workers
            </Typography>
            <IconButton onClick={toggleOpenAddWorker}>
              <ControlPoint color="primary"/>
            </IconButton>
          </div>
          <TableContainer component={Paper}>
            <Table className={classes.workerTable}>
              <TableHead>
                <TableRow>
                  <TableCell/>
                  <TableCell>IP</TableCell>
                  <TableCell>Status</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {workers.map((worker, i) => (
                  <WorkerTableRow
                    key={worker.ip}
                    toggleOpenAddImage={toggleOpenAddImage}
                    setSelectWorker={setSelectWorker}
                    worker={worker}
                    index={i}
                    workers={workers}
                    setWorkers={setWorkers}
                  />
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </div>
      </Paper>
    </ThemeProvider>
  );
}

ReactDOM.render(<App />, document.getElementById('root'));

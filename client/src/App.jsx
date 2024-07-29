import React from 'react';
import styled from 'styled-components';

const AppContainer = styled.div`
  text-align: center;
`;

function App() {
  return (
    <AppContainer>
      <h1>Welcome to LineUp</h1>
      <p>A custom Gomoku game</p>
    </AppContainer>
  );
}

export default App;
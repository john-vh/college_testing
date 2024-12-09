import React from 'react';
import { NonIdealState, Button, NonIdealStateIconSize, Icon, Colors } from '@blueprintjs/core';
import { LandingNavbar } from './LandingNavbar.tsx';

export const LoginSplash = ({ handleLogin, handleDialog }) => {
    return (
        <div>
            <LandingNavbar showHome={false} handleDialog={handleDialog} />

            <div
                style={{
                    height: 'calc(100vh - 64px)',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center'
                }}
            >
                <div style={{ textAlign: 'center', maxWidth: '600px' }}>
                    <Icon
                        icon="arrow-right"
                        size={60}
                        color={Colors.GOLD4}
                    />
                    <h2
                        className="bp4-heading"
                        style={{
                            marginBottom: '10px',
                            fontSize: '45px'
                        }}
                    >
                        Welcome to TestHive
                    </h2>

                    <p
                        className="bp4-text-large"
                        style={{
                            marginBottom: '25px',
                            color: '#5C7080' // Blueprint's default text color
                        }}
                    >
                        In order to use TestHive, please log in with a Google account
                    </p>

                    <Button
                        intent="primary"
                        style={{ background: Colors.GOLD4 }}
                        large={true}
                        icon="user"
                        onClick={handleLogin}
                    >
                        Sign in with Google
                    </Button>
                </div>
            </div>
        </div>
    );
};
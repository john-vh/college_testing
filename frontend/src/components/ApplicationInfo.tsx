import { Button, Card, CardList, Collapse, Colors, FormGroup, H3, H5, Icon, InputGroup, Label, TextArea } from "@blueprintjs/core";
import useAccountInfo from "../hooks/useAccountInfo.ts";
import React, { useEffect, useMemo, useState } from "react";
import { useApplicationInfo } from "../hooks/useApplicationInfo.ts";
import { useBusinessIds } from "../hooks/useBusinessIds.ts";
import { useAcceptApplication } from "../hooks/useAcceptApplication.ts";
import { useRejectApplication } from "../hooks/useRejectApplication.ts";
import { usePostingNames } from "../hooks/usePostingIds.ts";
import { useBusinessInfo } from "../hooks/useBusinessInfo.ts";

interface ApplicationInfoProps {
    isAdmin: boolean;
}

export const ApplicationInfo = ({ isAdmin }: ApplicationInfoProps) => {
    const [isOpen, setIsOpen] = useState<boolean[]>([]);
    const data = useApplicationInfo({ isAdmin });
    const names = usePostingNames(isAdmin);
    const acceptApplication = useAcceptApplication();
    const rejectApplication = useRejectApplication();

    const toggleCollapse = (index) => {
        setIsOpen((prevState) => ({
            ...prevState,
            [index]: !prevState[index],
        }));
    };

    const handleAcceptApplication = (entry, index) => {
        acceptApplication({ entry, index });
    }

    const handleRejectApplication = (entry, index) => {
        rejectApplication({ entry, index });
    }

    const filteredData = useMemo(() => {
        return data.filter((entry) => entry.applications.length > 0);
    }, [data]);

    return (
        <div className="content" style={{ width: '100%', maxWidth: '100%', overflow: 'hidden' }}>
            <CardList>
                {filteredData.map((entry, id) => (
                    <div style={{ display: "flex", flexDirection: "column" }}>
                        <Card key={id} style={{ width: '100%', maxWidth: '100%', overflow: 'hidden' }}>
                            <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between" }}>
                                <H5>{names.get(entry.post_id) || 'Unknown Posting'}</H5>
                                {entry.applications.length > 0 && <Button
                                    onClick={() => toggleCollapse(id)}
                                    rightIcon={isOpen[id] ? "chevron-up" : "chevron-down"}>
                                    {isOpen[id] ? 'Hide Applications' : 'Show Applications'}
                                </Button>}
                            </div>
                            <Collapse isOpen={isOpen[id]}>
                                <div style={{
                                    width: '100%',
                                    overflow: 'hidden',
                                    position: 'relative',
                                    padding: '20px 0'
                                }}>
                                    <div style={{
                                        display: "flex",
                                        flexDirection: "row",
                                        overflowX: "scroll",
                                        gap: "20px",
                                        padding: "10px",
                                    }}>
                                        {entry.applications.map((application, index) =>
                                            <Card interactive={true} >
                                                <H3>Application Information</H3>
                                                <FormGroup inline label="Posting"
                                                    labelFor="posting" >
                                                    <InputGroup id={"posting"} placeholder={entry.post_id.toString()} readOnly />
                                                </FormGroup>
                                                <FormGroup label="Name"
                                                    labelFor="name" >
                                                    <InputGroup id="name" placeholder={application.user.name} readOnly />
                                                </FormGroup>
                                                <FormGroup label="Email"
                                                    labelFor="email" >
                                                    <InputGroup id="email" placeholder={application.user.email} readOnly />
                                                </FormGroup>
                                                <Button onClick={() => handleAcceptApplication(entry, index)} style={{ background: Colors.VIOLET2, color: Colors.WHITE }}>Accept applicant</Button>
                                                <Button onClick={() => handleRejectApplication(entry, index)} style={{ margin: "10px" }}>Reject applicant</Button>
                                            </Card>
                                        )}
                                    </div>
                                </div>
                            </Collapse>
                        </Card>
                    </div>
                ))
                }
            </CardList >
        </div >);
}
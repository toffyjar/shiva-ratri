import * as React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { List, MD3Colors } from 'react-native-paper';

import { savedKundalisPageStyles as styles } from '../custom-styles/savedKundalisPageStyles';
import useKundliStore from '../store/kundliStore';

export default function SavedKundalis() {
  const { latitude, longitude, year, month, day, hour, minute, timeZone } = useKundliStore();
  const [kundalis, setKundalis] = React.useState([
    { id: 1, name: 'Kundali 1' },
    { id: 2, name: 'Kundali 2' },
    { id: 3, name: 'Kundali 3' },
  ]);

  const logInfo = () => {
    console.log('Latitude:', latitude);
    console.log('Longitude:', longitude);
    console.log('Date:', `${year}-${month}-${day}`);
    console.log('Time:', `${hour}:${minute}`);
    console.log('Time Zone:', timeZone);
  }

  return (
    <View style={styles.container}>
      <List.Section>
        {kundalis.map((kundali) => (
          <List.Item
            key={kundali.id}
            title={kundali.name}
            left={() => <List.Icon icon="file" />}
            onPress={() => {
              console.log('Item clicked!', kundali.id);
              logInfo();
            }}
          />
        ))}
      </List.Section>
    </View>
  );
}


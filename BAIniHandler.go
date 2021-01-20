package BAIniHandler

import "os"
import "bufio"
import "regexp"
import "errors"
import "strings"
import "strconv"
import "fmt"

type BAIniHandler interface {
	ReadString( section string, field string, defaultValue string ) ( string );
	ReadInteger( section string, field string, defaultValue int ) ( int );
	WriteString( section string, field string, value string );
	WriteInteger( section string, field string, value int ); 
	Save() ( error );
}

type baKeyValue map[string]interface{};

type BAIniHandle struct {
	fileName string;
	section map[string]baKeyValue
}

func NewBAIniHandler( fileName string, forceCreatingOnNotFound bool ) ( BAIniHandler, error ) {
	// I don´t like using defer. Forgive-me

	if _, err := os.Stat( fileName );( ( err != nil ) && ( forceCreatingOnNotFound ) || ( err == nil ) ) {
		var f *os.File;

		if f, err = os.Create( fileName ); ( err == nil ) {
			f.Close();
		} else {
			return nil, err;
		}
		
		if f, err = os.Open( fileName ); ( err != nil ) {
			return nil, err;
		} else {
			result := &BAIniHandle{ fileName: fileName, section: make( map[string]baKeyValue ) };
			if ( !readIniFile( f, result ) ) {
				f.Close();
				panic( errors.New( "Arquivo inconsistente!") );
			} else {
				f.Close();
			}
			return result, nil;
		}
	} else {
		return nil, err;
	}
}

func readIniFile( file *os.File, handler *BAIniHandle ) ( bool ) {
	var result bool = true;
	var section string;

	scanner := bufio.NewScanner(file);
	for scanner.Scan() {
		line := strings.TrimSpace( scanner.Text() );
		if hasRegex, err := regexp.MatchString( "(\\[{1})(\\w)+(\\]{1})", line ); ( err != nil ) {
			result = ( err != nil );
			if ( result ) { break; }
		} else {
			if ( hasRegex ) {
				section = strings.Replace( strings.Replace( line, "[", "", -1 ), "]", "", -1 );
			} else {
				if hasRegex, err := regexp.MatchString( "(\\w)+((\\s)+)?([=])((\\s)+)?(\\w)+", line ); ( ( err == nil ) && ( hasRegex ) ) {
					keyValue := strings.Split( line, "=" );
					if value, err := strconv.Atoi( strings.TrimSpace( keyValue[1] ) ); ( err == nil ) {
						handler.WriteInteger( section, strings.TrimSpace( keyValue[0] ), value );
					} else {
						handler.WriteString( section, strings.TrimSpace( keyValue[0] ), strings.TrimSpace( keyValue[1] ) );
					}
				} else {
					result = ( err != nil );
					if ( result ) { break; }
				}
			}
		}
	}	

	return result;
}

func ( ini * BAIniHandle ) Save()  ( error ) {
	if file, err := os.Create( ini.fileName ); ( err == nil ) {
		for section, content := range ini.section {			
			fmt.Fprintf( file, fmt.Sprintf ("[%v]\n", section ) );
			for key, value := range content {
				fmt.Fprintf( file, fmt.Sprintf ("%v=%v\n", key, value ) );
			}
		}
		file.Close();
		return nil;	
	} else {
		return err;
	}
}

func( ini *BAIniHandle) readSomething( section string, field string, defaultValue interface{} ) ( interface{} ) {
	keyValues := ini.section[ strings.ToUpper( section ) ];
	if ( keyValues != nil ) {
		result := keyValues[ strings.ToUpper( field ) ]; 
		if( result != nil ) {
			return result;
		}
	}
	return defaultValue;
}

func( ini *BAIniHandle) ReadString( section string, field string, defaultValue string ) ( string ) {
	return ini.readSomething( section, field, defaultValue ).(string);
}

func( ini *BAIniHandle) ReadInteger( section string, field string, defaultValue int ) ( int ) {
	return ini.readSomething( section, field, defaultValue ).(int);
}

func( ini *BAIniHandle) writeSomething( section string, field string, value interface{} ) {
	keyValues := ini.section[ strings.ToUpper( section ) ];
	if ( keyValues == nil ) {
		ini.section[ strings.ToUpper( section ) ] = make( baKeyValue );
		keyValues = ini.section[ strings.ToUpper( section ) ];
	}
	keyValues[ strings.ToUpper( field ) ] = value;
}

func( ini *BAIniHandle) WriteString( section string, field string, value string ) {
	ini.writeSomething( section, field, value );
}

func( ini *BAIniHandle) WriteInteger( section string, field string, value int ) {
	ini.writeSomething( section, field, value );
}
